package handlers

import (
	"net/http"
	"regexp"
	"slices"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	appmiddleware "coffee-site-server/internal/middleware"
)

// localDevOrigin matches any http(s)://localhost:<port> or 127.0.0.1:<port>
// origin. Vite picks the next free port when 5173 is taken, so pinning CORS
// to a single configured port is too brittle for local development.
var localDevOrigin = regexp.MustCompile(`^https?://(localhost|127\.0\.0\.1):\d+$`)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return slices.Contains(h.Cfg.AllowOrigins, origin) || localDevOrigin.MatchString(origin)
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-Guest-Id"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	requireAuth := appmiddleware.RequireAuth(h.Cfg.JWTSecret)
	optionalAuth := appmiddleware.OptionalAuth(h.Cfg.JWTSecret)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/verify-register", h.VerifyRegisterCode)
		r.Post("/login", h.Login)
		r.Post("/forgot-password", h.ForgotPassword)
		r.Post("/verify-reset-code", h.VerifyResetCode)
		r.Post("/reset-password", h.ResetPassword)

		r.Group(func(r chi.Router) {
			r.Use(requireAuth)
			r.Get("/me", h.Me)
			r.Post("/verify-permission", h.VerifyPermission)
		})
	})

	r.Route("/api/beans", func(r chi.Router) {
		r.Get("/", h.ListBeans)
		r.Get("/roast-levels", h.ListRoastLevels)
		r.Get("/process-methods", h.ListProcessMethods)
		r.Get("/{id}", h.GetBean)

		r.Group(func(r chi.Router) {
			r.Use(requireAuth, appmiddleware.RequireAdmin)
			r.Post("/", h.CreateBean)
			r.Put("/{id}", h.UpdateBean)
			r.Delete("/{id}", h.DeleteBean)
			r.Post("/roast-levels", h.AddRoastLevel)
			r.Patch("/roast-levels/{name}/toggle-hidden", h.ToggleRoastLevel)
			r.Post("/process-methods", h.AddProcessMethod)
			r.Patch("/process-methods/{name}/toggle-hidden", h.ToggleProcessMethod)
		})
	})

	r.Route("/api/utensils", func(r chi.Router) {
		r.Get("/", h.ListUtensils)
		r.Get("/categories", h.ListUtensilCategories)
		r.Get("/{id}", h.GetUtensil)

		r.Group(func(r chi.Router) {
			r.Use(requireAuth, appmiddleware.RequireAdmin)
			r.Post("/", h.CreateUtensil)
			r.Put("/{id}", h.UpdateUtensil)
			r.Delete("/{id}", h.DeleteUtensil)
			r.Post("/categories", h.AddUtensilCategory)
			r.Patch("/categories/{name}/toggle-hidden", h.ToggleUtensilCategory)
		})
	})

	r.Route("/api/cart", func(r chi.Router) {
		r.Get("/", h.ListCart)
		r.Post("/items", h.AddCartItem)
		r.Patch("/items/{productId}", h.UpdateCartItem)
		r.Delete("/items/{productId}", h.RemoveCartItem)
		r.Delete("/", h.ClearCart)
	})

	r.Route("/api/wishlist", func(r chi.Router) {
		r.Get("/", h.ListWishlist)
		r.Post("/{productId}/toggle", h.ToggleWishlistItem)
	})

	r.Route("/api/orders", func(r chi.Router) {
		r.Use(requireAuth)
		r.Get("/", h.ListOrders)
		r.Post("/checkout", h.Checkout)
	})

	r.Route("/api/contact", func(r chi.Router) {
		r.Use(optionalAuth)
		r.Post("/", h.CreateContactMessage)
	})

	return r
}
