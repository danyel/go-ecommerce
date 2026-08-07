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

func (h *shoppingBasketHandler) CreateShoppingBasket(w Http.ResponseWriter, _ *Http.Request) {
	sh, err := h.s.CreateShoppingBasket()
	if err != nil {
		CommonHandler.StatusInternalServerError(w)
		return
	}
	CommonHandler.WriteResponse(Http.StatusCreated, w, sh.ID)
}

func (h *shoppingBasketHandler) UpdateShoppingBasketItem(w Http.ResponseWriter, r *Http.Request) {
	var ai UpdateShoppingBasketItem
	var err error
	var shoppingBasket ShoppingBasket
	id, err := CommonHandler.GetId(r, "shoppingBasketId")
	if err != nil {
		CommonHandler.StatusBadRequest(w)
		return
	}

	if err = CommonHandler.ValidateRequest[UpdateShoppingBasketItem](r, &ai); err != nil {
		CommonHandler.StatusBadRequest(w)
		return
	}
	if shoppingBasket, err = h.s.UpdateShoppingBasketItem(id.ID, ai); err != nil {
		CommonHandler.StatusInternalServerError(w)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, shoppingBasket)
}

func (h *shoppingBasketHandler) GetShoppingBasket(w Http.ResponseWriter, r *Http.Request) {
	var err error
	var shoppingBasket ShoppingBasket
	id, err := CommonHandler.GetId(r, "shoppingBasketId")
	if err != nil {
		CommonHandler.StatusBadRequest(w)
		return
	}
	if shoppingBasket, err = h.s.GetShoppingBasket(id.ID); err != nil {
		CommonHandler.StatusInternalServerError(w)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, shoppingBasket)
}

func NewHandler(db *Database.DB, p Port.EventPublisher) ShoppingBasketHandler {
	s := NewService(db, p)
	return &shoppingBasketHandler{s}
}
