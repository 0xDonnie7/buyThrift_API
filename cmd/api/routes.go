package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/auth/signup", app.signupHandler)
		r.Post("/auth/login", app.loginHandler)
	})

	r.Route("/products", func(r chi.Router) {
		// public reads
		r.Get("/", app.getProducts)
		r.Get("/{id}", app.getProductByID)

		// admin-only writes
		r.Group(func(r chi.Router) {
			r.Use() //authentication middleware requireAuth requireAdmin

			r.Post("/", app.createProduct)
			r.Put("/{id}", app.updateProduct)
			r.Delete("/{id}", app.deleteProduct)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use() //authentication middleware

		r.Route("/cart", func(r chi.Router) {
			r.Get("/", app.getCartItem)
			r.Post("/items", app.addItemToCart)
			r.Patch("/items/{id}", app.updateItemQuantity)
			r.Delete("/items/{id}", app.deleteItemFromCart)
		})
	})

	return r
}

// PRODUCT ENDPOINT

// GET /products - list, with pagination, filtering, search
// GET /products/:id - detail
// POST /products - admin only
// PUT/PATCH /products/:id - admin only
// DELETE /products/:id - admin only (soft delete preferred, so old orders still reference valid data)

// CART ENDPOINT

//GET /cart - get current user's active cart
//POST /cart/items - add product (check stock!)
//PATCH /cart/items/:id - update quantity
//DELETE /cart/items/:id -remove item

// Note: server-side price calculation only, never trust a price sent from client
