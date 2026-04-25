package meals

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{db: db}
}

func (r *MealRepository) ListAllMeals(ctx context.Context) ([]*Meal, error) {
	return generic_queries.List[*Meal](ctx, r.db)
}

func (r *MealRepository) ListMealsAndFoodsAndIngredients(ctx context.Context) ([]*Meal, error) {
	return gorm.G[*Meal](r.db).
		Scopes(wherePiid(ctx)).
		Preload("Foods", nil).
		Preload("Foods.Ingredient", nil).
		Find(ctx)
}

func (r *MealRepository) GetMealAndFoods(ctx context.Context, id uint) (*Meal, error) {
	return gorm.G[*Meal](r.db).
		Scopes(wherePiid(ctx)).
		Preload("Foods", nil).
		Where("id = ?", id).
		First(ctx)
}

func (r *MealRepository) CreateMeal(ctx context.Context, meal *Meal) (uint, error) {
	err := generic_queries.Create(ctx, r.db, meal)
	return meal.ID, err
}

func (r *MealRepository) DeleteMeal(ctx context.Context, id uint) error {
	return generic_queries.Delete[*Meal](ctx, r.db, id)
}

func (r *MealRepository) UpdateMeal(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "meals", id, values)
}

func (r *MealRepository) ListFoodsByMeal(ctx context.Context, mealId uint) ([]Food, error) {
	return gorm.G[Food](r.db).
		Preload("Ingredient", nil).
		Where("meal_id = ?", mealId).
		Find(ctx)
}

func (r *MealRepository) ListFoodsByIds(ctx context.Context, ids []uint) ([]Food, error) {
	return gorm.G[Food](r.db).Where("id IN ?", ids).Find(ctx)
}

func (r *MealRepository) CreateFoodAndIngredient(ctx context.Context, food *Food, ingredientName string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var ingredient = Ingredient{Name: ingredientName}
		if err := generic_queries.Create(ctx, tx, &ingredient); err != nil {
			return err
		}
		food.IngredientID = ingredient.ID
		return tx.Clauses(clause.Returning{}).Create(food).Error
	})
}

func (r *MealRepository) CreateFood(ctx context.Context, food *Food) error {
	return gorm.G[Food](r.db, clause.Returning{}).Create(ctx, food)
}

func (r *MealRepository) DeleteFood(ctx context.Context, id uint) error {
	rows, err := gorm.G[Food](r.db, clause.Returning{}).
		Where("id = ?", id).
		Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *MealRepository) FirstFood(ctx context.Context, id uint) (Food, error) {
	return gorm.G[Food](r.db).Where("id = ?", id).First(ctx)
}

func (r *MealRepository) UpdateFood(ctx context.Context, id uint, values map[string]any) error {
	rows, err := gorm.G[map[string]any](r.db).
		Table("foods").
		Where("id = ?", id).
		Updates(ctx, values)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *MealRepository) CreateFoods(ctx context.Context, foods []Food) error {
	err := gorm.G[Food](r.db).CreateInBatches(ctx, &foods, 100)
	return err
}
