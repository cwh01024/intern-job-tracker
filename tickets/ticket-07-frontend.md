# Ticket 07: Frontend Dashboard

## Context
A simple dashboard to view discovered jobs, notification history, and manage settings. Built with vanilla HTML/CSS/JS for simplicity.

## Goals
1. Create responsive dashboard UI
2. Display jobs in sortable table
3. Show notification history
4. Add settings form for recipient config
5. Manual refresh button

## Dependencies
- Ticket 06 (API must be ready)

## Acceptance Criteria
- [ ] `web/index.html` - main page structure
- [ ] `web/style.css` - responsive styling
- [ ] `web/app.js` - API integration
- [ ] Dashboard loads jobs from API
- [ ] Settings can be updated

## Technical Details

### Pages/Sections
1. **Jobs Table**: Company, Title, Date, Link
2. **Notifications**: History with timestamps
3. **Settings**: Recipient phone, schedule time
4. **Actions**: Refresh button

### UI Features
- Dark mode support
- Mobile responsive
- Loading states
- Error messages

### Files to Create
```
web/index.html
web/style.css
web/app.js
```

## Implementation Notes
- Use fetch() for API calls
- Use CSS Grid/Flexbox for layout
- Add favicon and meta tags

## Estimated Tokens
~400 tokens (frontend is simpler)
