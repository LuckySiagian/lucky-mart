package services

import (
	"strings"
	"testing"

	"backend-sembako/internal/constants"
	"backend-sembako/internal/models"
)

func TestFormatRupiah(t *testing.T) {
	cases := map[float64]string{0: "Rp0", 950: "Rp950", 15000: "Rp15.000", 1234567: "Rp1.234.567"}
	for in, want := range cases {
		if got := formatRupiah(in); got != want {
			t.Errorf("formatRupiah(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestRenderTemplateBodyDropsUnknownPlaceholders(t *testing.T) {
	got := renderTemplateBody("Halo {{customer_name}}, pesanan {{order_number}} {{missing}}siap.", map[string]interface{}{
		"customer_name": "Sari",
		"order_number":  "ORD-1",
	})
	want := "Halo Sari, pesanan ORD-1 siap."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNewReferenceFormat(t *testing.T) {
	ref := newReference("ORD")
	parts := strings.Split(ref, "-")
	if len(parts) != 3 || parts[0] != "ORD" || len(parts[1]) != 8 || len(parts[2]) != 6 {
		t.Errorf("unexpected reference %q", ref)
	}
	if newReference("ORD") == ref {
		t.Error("references should differ")
	}
}

func TestWholesalePriceOnlyForB2B(t *testing.T) {
	unit := &models.ProductUnit{RetailPrice: 10000, WholesalePrice: 9000, MinWholesaleQty: 5}
	if p := calculateUnitPrice(unit, constants.RoleB2CCustomer, 10); p != 10000 {
		t.Errorf("retail customer got %v", p)
	}
	if p := calculateUnitPrice(unit, constants.RoleB2BCustomer, 4); p != 10000 {
		t.Errorf("B2B below minimum got %v", p)
	}
	if p := calculateUnitPrice(unit, constants.RoleB2BCustomer, 5); p != 9000 {
		t.Errorf("B2B at minimum got %v", p)
	}
}

func TestCanManage(t *testing.T) {
	if canManage(constants.RoleStoreManager, constants.RoleSuperAdmin) {
		t.Error("store manager must not manage a super admin")
	}
	if canManage(constants.RoleStoreManager, constants.RoleStoreManager) {
		t.Error("store manager must not manage another manager")
	}
	if !canManage(constants.RoleStoreManager, constants.RoleDriver) {
		t.Error("store manager should manage drivers")
	}
	if !canManage(constants.RoleSuperAdmin, constants.RoleSuperAdmin) {
		t.Error("super admin should manage everyone")
	}
}
