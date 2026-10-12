package constants

import "testing"

func TestOrderTransitions(t *testing.T) {
	cases := []struct {
		from, to OrderStatus
		ok       bool
	}{
		{OrderStatusConfirmed, OrderStatusProcessingPacking, true},
		{OrderStatusProcessingPacking, OrderStatusOutForDelivery, true},
		{OrderStatusOutForDelivery, OrderStatusCompleted, true},
		{OrderStatusConfirmed, OrderStatusCancelled, true},
		// Skipping ahead is allowed (e.g. picked up at the store right away).
		{OrderStatusConfirmed, OrderStatusCompleted, true},
		{OrderStatusProcessingPacking, OrderStatusCompleted, true},
		// A customer who never takes the goods can still be cancelled.
		{OrderStatusOutForDelivery, OrderStatusCancelled, true},
		// Going backwards or leaving a final status is refused.
		{OrderStatusOutForDelivery, OrderStatusProcessingPacking, false},
		{OrderStatusCompleted, OrderStatusConfirmed, false},
		{OrderStatusCompleted, OrderStatusCancelled, false},
		{OrderStatusCancelled, OrderStatusConfirmed, false},
		// A transfer order is confirmed only by approving its payment.
		{OrderStatusPendingPayment, OrderStatusConfirmed, false},
		{OrderStatusPendingPayment, OrderStatusCompleted, false},
	}
	for _, c := range cases {
		if got := c.from.CanTransitionTo(c.to); got != c.ok {
			t.Errorf("%s -> %s: got %v, want %v", c.from, c.to, got, c.ok)
		}
	}
}

func TestRoleHelpers(t *testing.T) {
	if !RoleDriver.IsStaff() || RoleB2BCustomer.IsStaff() {
		t.Error("IsStaff misclassifies roles")
	}
	if !RoleStoreManager.IsManager() || RoleWarehouseStaff.IsManager() {
		t.Error("IsManager misclassifies roles")
	}
}
