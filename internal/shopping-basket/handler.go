package shopping_basket

import (
	Http "net/http"

	CommonHandler "github.com/danyel/ecommerce/internal/common/handler"
	Port "github.com/danyel/ecommerce/internal/common/port"
	Database "gorm.io/gorm"
)

type ShoppingBasketHandler interface {
	CreateShoppingBasket(w Http.ResponseWriter, r *Http.Request)
	UpdateShoppingBasketItem(w Http.ResponseWriter, r *Http.Request)
	GetShoppingBasket(w Http.ResponseWriter, r *Http.Request)
}

type shoppingBasketHandler struct {
	s ShoppingBasketService
}

func (h *shoppingBasketHandler) CreateShoppingBasket(w Http.ResponseWriter, r *Http.Request) {
	sh, err := h.s.CreateShoppingBasket()
	if err != nil {
		CommonHandler.StatusInternalServerError(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusCreated, w, r, sh.ID)
}

// UpdateShoppingBasketItem web handler function that will update the shopping basket
func (h *shoppingBasketHandler) UpdateShoppingBasketItem(w Http.ResponseWriter, r *Http.Request) {
	var updateShoppingBasketItem UpdateShoppingBasketItem
	var err error
	var shoppingBasket ShoppingBasket
	shoppingBasketId, err := CommonHandler.GetId(r, "shoppingBasketId")
	if err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}

	if err = CommonHandler.ValidateRequest[UpdateShoppingBasketItem](r, &updateShoppingBasketItem); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}
	if shoppingBasket, err = h.s.UpdateShoppingBasketItem(shoppingBasketId.ID, updateShoppingBasketItem); err != nil {
		CommonHandler.StatusInternalServerError(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, r, shoppingBasket)
}

func (h *shoppingBasketHandler) GetShoppingBasket(w Http.ResponseWriter, r *Http.Request) {
	var err error
	var shoppingBasket ShoppingBasket
	id, err := CommonHandler.GetId(r, "shoppingBasketId")
	if err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}
	if shoppingBasket, err = h.s.GetShoppingBasket(id.ID); err != nil {
		CommonHandler.StatusInternalServerError(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, r, shoppingBasket)
}

func NewHandler(db *Database.DB, p Port.EventPublisher) ShoppingBasketHandler {
	s := NewService(db, p)
	return &shoppingBasketHandler{s}
}
