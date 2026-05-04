package service_test

import (
	"context"
	"testing"

	"github.com/archera/shipping-service-ver2/internal/service"
)

func TestAutoDispatchPickUp_EmptyOrderID(t *testing.T) {
	
	dispatchSvc := service.NewDispatchService(nil, nil, nil)
	ctx := context.Background()

	// Act: Panggil fungsi dengan orderID kosong
	_, err := dispatchSvc.AutoDispatchPickUp(ctx, "", "Bandung")

	if err == nil {
		t.Fatalf("Ekspektasi error karena OrderID kosong, tapi malah sukses")
	}

	expectedErrMsg := "order ID tidak boleh kosong"
	if err.Error() != expectedErrMsg {
		t.Errorf("Ekspektasi error '%s', tapi dapat '%s'", expectedErrMsg, err.Error())
	}
}