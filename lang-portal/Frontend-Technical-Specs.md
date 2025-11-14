# Frontend Technical Spec

## Pages 

### Dashboard `/Dashboard`

#### Purpose
The purpose of this page is to provide a summary of learning and act as a default page when the user visits the web-app 

#### Components
- Last Study Session
  - show last activity used and time taken
  - summarises wrong vs correct from last activity
  - has a link to the group
- Study Progress
  - total words studied- eg. 3/124
      across all study Session show the total words studied out of all possible words in our database
  - Display a mastery progress- eg. 0% 
- Quick Stats
  - success rate - eg. 80%
  - total study session - eg. 4
  - total active groups - eg. 3
  - study streak - eg. 2 days
-Start Studying Button
  - goes to study activities page

#### Needed API endpoints
- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats

### Study Activities Index `/study_activities`

#### Purpose
The purpose of this page is to show a collection of study activities with a thumbnail and its name, to either launch or view the study activity.

#### Components
- Study Activity Card
  - show thumbnail of the study activity
  - the name of the study activity
  - a launch button to take us to the launch page
  - the view page to view more information about past study sessions for this study activity
 
 #### Needed API endpoints
- GET /api/study_activities

### Study Activity show `/study_activities/:id`

#### Purpose
The purposeof this page is to show the details of specific study activity and its past study sessions.

#### Components
- Name of study activity
- Thumbnail of study activity
- Description of study activity
- Launch Button
- Study Activities paginated List
  - id
  - activity name
  - group name
  - start time
  - end time(infered by the last word_review_item submitted)
  - number of review items

#### Needed API endpoints
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions

### Study Activity Launch `/study_activities/:id/launch`

#### Purpose
The purpose of this page is to launch a study activity.

#### Components
- Name of study activity
- Launch from
  - select field for word group
  - launch now Button

## Behaviour
-After the form is submitted a new tab opens with the study activity based on its URL provided in its database

-also after the form is submitted the page will redirect to the study session show Page

#### Needed API endpoints
- POST /api/study_activities

### Words Index `/words`

#### Purpose
The purpose of this page is to show all words in the database.

#### Components
- paginated Word List
  - Columns
    - japenese
    - romaji
    - english
    - Correct Count
    - wrong Count
  - Pagination with 100 words per page
  - Clicking the japenese word will take us to the word show page
  
#### Needed API endpoints
- GET /api/words

### Word show `/words/:id`

#### Purpose
The purpose of this page is to show the details of a specific word.

#### Components
- japenese
- romaji
- english
- study statistics
  - correct count
  - wrong count
- Word groups
  - shown a series of pills eg. tags
  - when the group name is clicked it will take us to the group show page

  
#### Needed API endpoints
- GET /api/words/:id

### Word Group Index `/word_groups`

#### Purpose
The purpose of this page is to show a list of groups in the database.

#### Components
- paginated Word Group List
  - Columns
    - group name
    - number of words
  - Clicking the name of the word group will take us to the word group show page

#### Needed API endpoints
- GET /api/groups

### Group show `/groups/:id`

#### Purpose
The purpose of this page is to show the details of a specific group.

#### Components
- group name
- group statistics
  - total word count
- word in group (Paginated list of words)
  - should use the same component as the word index page
- study sessions (Paginated list of study sessions)
  - should use the same component as the study sessions index page 

#### Needed API endpoints
- GET /api/groups/:id(the name of the group stats)
- GET /api/groups/:id/words
- GET /api/groups/:id/study_sessions

### Study Sessions Index `/study_sessions`

#### Purpose
The purpose of this page is to show list of study sessions in the database.

#### Components
- paginated Study Session List
  - Columns
    - id
    - activity name
    - group name
    - start time
    - end time(infered by the last word_review_item submitted)
    - number of review items
  - Clicking the Study Session id will take us to the study session show page

#### Needed API endpoints
- GET /api/study_sessions

### Study Session show `/study_sessions/:id`

#### Purpose
The purpose of this page is to show the details of a specific study session.

#### Components
- study session details
  - activity name
  - group name
  - start time
  - end time(infered by the last word_review_item submitted)
  - number of review items
- words review items (Paginated list of words)
  - should use the same component as the review items index page

#### Needed API endpoints
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/words

### Settings `/settings`

#### Purpose
The purpose of this page is to make configuration to the study portal

#### Components
  - Theme selection eg. light, dark, system default
  - Language selection eg. english, japanese, korean, chinese
  - Reset History Button
    - This will delete all study sessions and word review items
  -full reset Button
    - This will drop all tables and recreate with seed data

#### Needed API endpoints
- POST /api/reset_history
- POST /api/full_reset
