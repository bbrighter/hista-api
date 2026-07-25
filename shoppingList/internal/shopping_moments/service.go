package shoppingmoments

import (
	"context"

	appError "encore.app/errors"
	"encore.app/shoppingList/internal/moments"
	sl "encore.app/shoppingList/internal/shoppingList"
	"gorm.io/gorm"
)

type ShoppingMomentsService struct {
	uow *UnitOfWork
	mom *moments.MomentRepo
}

func NewShoppingMomentsService(db *gorm.DB) *ShoppingMomentsService {
	uow := NewUnitOfWork(db)
	mom := moments.NewMomentRepo()
	return &ShoppingMomentsService{uow: uow, mom: mom}
}

func (s *ShoppingMomentsService) List() *sl.ShoppingListRepo {
	return s.uow.ShoppingList()
}

func (s *ShoppingMomentsService) CreateOrFirstList(ctx context.Context) (sl.List, error) {
	list, err := s.List().FirstShoppingList(ctx)
	if err == nil {
		return list, nil
	}
	if err != gorm.ErrRecordNotFound {
		return sl.List{}, err
	}

	newList := sl.List{}
	if err = s.List().CreateShoppingList(ctx, &newList); err != nil {
		return newList, err
	}
	err = s.mom.UpdateItems(ctx)
	return newList, err
}

func (s *ShoppingMomentsService) DeleteList(ctx context.Context, id uint) error {
	err := s.uow.WithTransaction(ctx, func(uow *UnitOfWork) error {
		items, err := uow.ShoppingList().ListItemsForList(ctx, id)
		if err != nil {
			return err
		}

		for _, item := range items {
			if !item.Checked {
				return appError.ErrCannotDelete
			}
		}

		if err := uow.ShoppingList().DeleteShoppingList(ctx, id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return s.mom.UpdateItems(ctx)
}

func (s *ShoppingMomentsService) ForceDeleteList(ctx context.Context, id uint) error {
	if err := s.List().DeleteShoppingList(ctx, id); err != nil {
		return err
	}
	return s.mom.UpdateItems(ctx)
}

func (s *ShoppingMomentsService) DeleteListCreateNewAndMoveItems(ctx context.Context, id uint) (sl.List, error) {
	var uncheckedItems []sl.Item
	var newList = &sl.List{}

	err := s.uow.WithTransaction(ctx, func(uow *UnitOfWork) error {
		items, err := uow.ShoppingList().ListItemsForList(ctx, id)
		if err != nil {
			return err
		}

		var uncheckedItemIds []uint
		for i := range items {
			if !items[i].Checked {
				items[i].ListId = id
				uncheckedItemIds = append(uncheckedItemIds, items[i].ID)
				uncheckedItems = append(uncheckedItems, items[i])
			}
		}

		if err := uow.ShoppingList().DeleteShoppingList(ctx, id); err != nil {
			return err
		}

		if err := uow.ShoppingList().CreateShoppingList(ctx, newList); err != nil {
			return err
		}

		if len(uncheckedItemIds) > 0 {
			if err := uow.ShoppingList().UpdateItemsList(ctx, uncheckedItemIds, newList.ID); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return sl.List{}, err
	}

	newList.Items = uncheckedItems
	return *newList, s.mom.UpdateItems(ctx)
}

func (s *ShoppingMomentsService) UpdateProduct(ctx context.Context, id uint, name *string, archive *bool) error {
	values := make(map[string]any)
	if name != nil {
		values["name"] = *name
	}
	if archive != nil {
		values["archived"] = *archive
	}
	if err := s.List().UpdateProduct(ctx, id, values); err != nil {
		return err
	}
	return s.mom.UpdateProducts(ctx)
}

func (s *ShoppingMomentsService) DeleteProduct(ctx context.Context, id uint) error {
	if err := s.List().DeleteProduct(ctx, id); err != nil {
		return err
	}
	return s.mom.UpdateProducts(ctx)
}

func (s *ShoppingMomentsService) AddItemByProductId(ctx context.Context, listId uint, productId uint) (uint, error) {
	var item = sl.Item{ProductId: productId, ListId: listId}
	if err := s.List().CreateItem(ctx, &item); err != nil {
		return 0, err
	}
	return item.ID, s.mom.UpdateItems(ctx)
}

func (s *ShoppingMomentsService) AddItemByName(ctx context.Context, listId uint, name string) (sl.Item, sl.Product, error) {
	var product = sl.Product{Name: name}
	var item = sl.Item{ListId: listId}

	err := s.uow.WithTransaction(ctx, func(uow *UnitOfWork) error {
		err := uow.ShoppingList().UpsertProduct(ctx, &product)
		if err != nil {
			return err
		}
		item.ProductId = product.ID
		return uow.ShoppingList().CreateItem(ctx, &item)
	})
	if err != nil {
		return item, product, err
	}
	if err := s.mom.UpdateItems(ctx); err != nil {
		return item, product, err
	}
	if err := s.mom.UpdateProducts(ctx); err != nil {
		return item, product, err
	}
	return item, product, nil
}

func (s *ShoppingMomentsService) DeleteItems(ctx context.Context, ids []uint) error {
	if err := s.List().DeleteItems(ctx, ids); err != nil {
		return err
	}
	return s.mom.UpdateItems(ctx)
}

func (s *ShoppingMomentsService) UpdateItem(ctx context.Context, id uint, quantity *uint8, checked *bool) error {
	values := make(map[string]any)
	if quantity != nil {
		values["quantity"] = *quantity
	}
	if quantity != nil && *quantity == 0 { // 0 should be treated as nil
		values["quantity"] = nil
	}
	if checked != nil {
		values["checked"] = *checked
	}
	if err := s.List().UpdateItem(ctx, id, values); err != nil {
		return err
	}
	return s.mom.UpdateItems(ctx)
}

func (s *ShoppingMomentsService) GetMoment(ctx context.Context) (moments.Moment, error) {
	return s.mom.GetMoments(ctx)
}

func (s *ShoppingMomentsService) ListProducts(ctx context.Context) ([]*sl.Product, error) {
	return s.List().ListProducts(ctx)
}
