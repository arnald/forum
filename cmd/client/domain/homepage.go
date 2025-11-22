package domain

// Data for Homepage.
type HomePageData struct {
	ActivePage string
	Categories []Category
	User       *LoggedInUser
}
