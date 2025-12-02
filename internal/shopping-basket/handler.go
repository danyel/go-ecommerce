package shopping_basket

import (
	"net/http"

	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"gorm.io/gorm"
)

type ShoppingBasketHandler interface {
	CreateShoppingBasket(w http.ResponseWriter, r *http.Request)
	UpdateShoppingBasketItem(w http.ResponseWriter, r *http.Request)
	GetShoppingBasket(w http.ResponseWriter, r *http.Request)
}

type shoppingBasketHandler struct {
	s ShoppingBasketService
	h commonHandler.ResponseHandler
}

func (h *shoppingBasketHandler) CreateShoppingBasket(w http.ResponseWriter, _ *http.Request) {
	sh, err := h.s.CreateShoppingBasket()
	if err != nil {
		h.h.StatusInternalServerError(w)
		return
	}

	h.h.WriteResponse(http.StatusCreated, w, ShoppingId{ID: sh.ID})
}

func (h *shoppingBasketHandler) UpdateShoppingBasketItem(w http.ResponseWriter, r *http.Request) {
	var ai UpdateShoppingBasketItem
	var err error
	var shoppingBasket ShoppingBasket
	id, err := commonHandler.GetId(r, "shoppingBasketId")
	if err != nil {
		h.h.StatusBadRequest(w)
		return
	}

	if err = commonHandler.ValidateRequest[UpdateShoppingBasketItem](r, &ai); err != nil {
		h.h.StatusBadRequest(w)
		return
	}
	if shoppingBasket, err = h.s.UpdateShoppingBasketItem(id, ai); err != nil {
		h.h.StatusInternalServerError(w)
		return
	}
	h.h.WriteResponse(http.StatusOK, w, shoppingBasket)
}

func (h *shoppingBasketHandler) GetShoppingBasket(w http.ResponseWriter, r *http.Request) {
	var err error
	var shoppingBasket ShoppingBasket
	id, err := commonHandler.GetId(r, "shoppingBasketId")
	if err != nil {
		h.h.StatusBadRequest(w)
		return
	}
	if shoppingBasket, err = h.s.GetShoppingBasket(id); err != nil {
		h.h.StatusInternalServerError(w)
		return
	}
	h.h.WriteResponse(http.StatusOK, w, shoppingBasket)
}

func NewHandler(db *gorm.DB) ShoppingBasketHandler {
	s := NewService(db)
	return &shoppingBasketHandler{s, commonHandler.NewResponseHandler()}
}
