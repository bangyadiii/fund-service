package campaign

type CreateCampaignInput struct {
	UserID           uint32 `json:"user_id" binding:"required" validate:"numeric"`
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description" binding:"required"`
	Perks            string `json:"perks"`
	BackerCount      uint   `json:"backer_count" binding:"required,numeric"`
	GoalAmount       uint   `json:"goal_amount" binding:"required,numeric"`
	Slug             string `json:"slug" binding:"required"`
}
