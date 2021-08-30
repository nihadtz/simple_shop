package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/nihadtz/simple_shop/models"
	"github.com/nihadtz/simple_shop/services"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
)

type Payments struct {
}

func (pc Payments) ViaStripe(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user := context.Get(req, "user").(models.User)
	purchaseID, _ := strconv.Atoi(params.ByName("id"))

	var purchase models.Purchase

	err := purchase.Get(purchaseID)

	if err != nil {
		services.LogError("Error reading Purchase", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	if user.ID != purchase.UserID {
		err := errors.New("you are not allowed to view this purchase")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusForbidden, err.Error())
		return
	}

	if purchase.Status == "paid" {
		err := errors.New("purchase has been already paid")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusNotFound, err.Error())
		return
	}

	stripe.Key = services.Configuration.StripeSK

	stripeCharge, err := charge.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(int64(purchase.Total) * 100),
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")}, // this will be sent from frontend, used like this for testing purpose
		ReceiptEmail: stripe.String(user.Email),
	})

	var payment = models.Payment{
		Date:       time.Now().Unix(),
		PurchaseID: purchase.ID,
		Gateway:    "stripe",
	}

	if err != nil {
		services.LogError("Error payment via Stripe", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())

		payment.Status = "failed"
		payment.Error = err.Error()
		payment.GatewayPaymentID = stripeCharge.ID

		err = payment.Create(purchase)

		if err != nil {
			services.LogError("Error creating Payment", err)
			services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
		}

		return
	}

	payment.Status = "ok"
	payment.GatewayPaymentID = stripeCharge.ID

	err = payment.Create(purchase)

	if err != nil {
		services.LogError("Error creating Payment", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
		return
	}

	services.Renderer.Render(res, http.StatusOK, payment)
}
