#TODO:

 - UI
   - DONE Create favicon
   - DONE Create landing image
   - DONE Embed static/* into binary
   - DONE Convert to Alpine.js
   - DONE Convert to tailwind
   - DONE List projects separately
   - DONE Edit environment
   - DONE Edit project
   - Validation for forms
   - DONE Show/hide deleted
 - API
   - DONE IN PROGRESS Create CRUD API endpoints for projects, environments
   - DONE Allow CRUD actions from UI for projects, environments
   - DONE User better routing to handle query params: https://github.com/benhoyt/go-routing/blob/master/retable/route.go
 - Write tests
 - Address all linter issues
 - Upgrade deprecation warning in serve.go
 - Fix README so it looks presentable on github
 - Check for os correctly - store db appropriately
 - Create basic config if not exists (and not already present in home dir)
 - Use config for:
   - setting app_env
   - Storing config from template (if not exists)
   - Storing db
 - Show application version on console and web ui
 - Break down all styles into classes in styles.css
 - Fix tz for environment dates
 - Allow specifying custom:
   - port
   - host
   - db (for development)
