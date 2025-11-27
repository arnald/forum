package domain

// Data for Homepage.
type HomePageData struct {
	User       *LoggedInUser
	ActivePage string
	Categories []Category
}
