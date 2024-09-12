package test

const MIN_RESUME string = `{
  "name": "name",
  "email": "email",
  "phoneNumber": "phoneNumber",
  "prelude": "prelude"
}`

const FULL_RESUME string = `{
  "name": "name",
  "email": "email",
  "phoneNumber": "phoneNumber",
  "prelude": "prelude",
  "location": "location",
  "linkedIn": "linkedIn",
  "github": "github",
  "facebook": "facebook",
  "instagram": "instagram",
  "twitter": "twitter",
  "portfolio": "portfolio"
}`

const MIN_EDUCATION string = `{
  "degreeType": "degreeType",
  "fieldOfStudy": "fieldOfStudy",
  "institution": "institution",
  "began": "1970-01-01T00:00:00.000Z",
  "current": true,
  "gpa": "string"
}`

const FULL_EDUCATION string = `{
  "degreeType": "degreeType",
  "fieldOfStudy": "fieldOfStudy",
  "institution": "institution",
  "began": "1970-01-01T00:00:00.000Z",
  "current": true,
  "location": "location",
  "finished": "1970-01-01T00:00:00.000Z",
  "gpa": "string"
}`

const MIN_WORK_EXPERIENCE string = `{
  "employer": "employer",
  "title": "title",
  "began": "1970-01-01T00:00:00.000Z",
  "current": true,
  "responsibilities": [
    "string"
  ]
}`

const FULL_WORK_EXPERIENCE string = `{
  "employer": "employer",
  "title": "title",
  "began": "1970-01-01T00:00:00.000Z",
  "current": true,
  "location": "location",
  "finished": "1970-01-01T00:00:00.000Z",
  "responsibilities": [
    "string"
  ]
}`

const PROJECT string = `{
  "name": "name",
  "description": "description",
  "role": "role",
  "responsibilities": [
    "responsibility"
  ]
}`
