import express from 'express'
import compression from 'compression'
import morgan from 'morgan'
import cors from 'cors'
import http from 'http'
import path from 'path'
import yaml from 'yamljs'
import * as api from './api'
import { accessTokenAuth } from './security'
import { connector, summarise } from 'swagger-routes-express'
import * as users from './users'
import * as resumeRepos from './resume-repos'
import * as resumeVersions from './resume-versions'
import * as resumeApplications from './resume-applications'
import * as interviewProgress from './interview-progress'
import * as jobPostings from './job-postings'
import * as interviewQuestions from './interview-questions'

const app = express()
const port = 9528

// Compression
app.use(compression())
// Logger
app.use(morgan('dev'))
// Enable CORS
app.use(cors())
// POST, PUT, DELETE body parser
app.use(express.json({ limit: '20mb' }))
app.use(express.urlencoded({
  limit: '20mb',
  extended: false
}))
// No cache
app.use((req, res, next) => {
  res.header('Cache-Control', 'private, no-cache, no-store, must-revalidate')
  res.header('Pragma', 'no-cache')
  res.header('Expires', '-1')
  next()
})

// ========== 简历管理系统的路由配置 ==========
// 用户相关
app.post('/api/user/register', users.register)
app.post('/api/user/email-code', users.sendEmailCode)
app.post('/api/user/login', users.login)
app.post('/api/user/logout', users.logout)
app.get('/api/user/profile', users.getUserProfile)
app.put('/api/user/profile', users.updateUserProfile)

// 简历仓库相关
app.get('/api/resume-repos', resumeRepos.getResumeRepos)
app.get('/api/resume-repos/:id', resumeRepos.getResumeRepoById)
app.post('/api/resume-repos', resumeRepos.createResumeRepo)
app.put('/api/resume-repos/:id', resumeRepos.updateResumeRepo)
app.delete('/api/resume-repos/:id', resumeRepos.deleteResumeRepo)

// 简历版本相关
app.get('/api/resume-versions', resumeVersions.getResumeVersions)
app.get('/api/resume-versions/:id', resumeVersions.getResumeVersionById)
app.post('/api/resume-versions', resumeVersions.createResumeVersion)
app.put('/api/resume-versions/:id', resumeVersions.updateResumeVersion)
app.delete('/api/resume-versions/:id', resumeVersions.deleteResumeVersion)
app.post('/api/resume-versions/:id/activate', resumeVersions.activateResumeVersion)

// 投递记录相关
app.get('/api/resume-applications', resumeApplications.getResumeApplications)
app.get('/api/resume-applications/:id', resumeApplications.getResumeApplicationById)
app.post('/api/resume-applications', resumeApplications.createResumeApplication)
app.put('/api/resume-applications/:id', resumeApplications.updateResumeApplication)
app.delete('/api/resume-applications/:id', resumeApplications.deleteResumeApplication)

// 岗位库相关
app.get('/api/job-postings', jobPostings.getJobPostings)
app.get('/api/job-postings/:id', jobPostings.getJobPostingById)
app.post('/api/job-postings', jobPostings.createJobPosting)
app.put('/api/job-postings/:id', jobPostings.updateJobPosting)
app.delete('/api/job-postings/:id', jobPostings.deleteJobPosting)
app.post('/api/job-postings/parse', jobPostings.parseJobPosting)

// 面试进展相关
app.get('/api/interview-progress', interviewProgress.getInterviewProgressList)
app.post('/api/interview-progress', interviewProgress.createInterviewProgress)
app.put('/api/interview-progress/:id', interviewProgress.updateInterviewProgress)
app.delete('/api/interview-progress/:id', interviewProgress.deleteInterviewProgress)

// 面试题库相关
app.get('/api/interview-questions', interviewQuestions.getInterviewQuestions)
app.get('/api/interview-questions/:id', interviewQuestions.getInterviewQuestionById)
app.post('/api/interview-questions', interviewQuestions.createInterviewQuestion)
app.put('/api/interview-questions/:id', interviewQuestions.updateInterviewQuestion)
app.delete('/api/interview-questions/:id', interviewQuestions.deleteInterviewQuestion)
app.put('/api/interview-questions/:id/favorite', interviewQuestions.updateInterviewQuestionFavorite)
app.get('/api/question-occurrences', interviewQuestions.getQuestionOccurrences)
app.post('/api/question-occurrences', interviewQuestions.createQuestionOccurrence)
app.delete('/api/question-occurrences/:id', interviewQuestions.deleteQuestionOccurrence)
app.get('/api/flashcards', interviewQuestions.getFlashcards)
app.post('/api/flashcard-reviews', interviewQuestions.createFlashcardReview)

// Read and swagger config file
const apiDefinition = yaml.load(path.resolve(__dirname, 'swagger.yml'))
// Create mock functions based on swaggerConfig
const options = {
  security: {
    AccessTokenAuth: accessTokenAuth
  }
}
const connectSwagger = connector(api, apiDefinition, options)
connectSwagger(app)
// Print swagger router api summary
const apiSummary = summarise(apiDefinition)
console.log(apiSummary)

// Catch 404 error
app.use((req, res) => {
  const err = new Error('Not Found')
  res.status(404).json({
    message: err.message,
    error: err
  })
})

// Create HTTP server.
const server = http.createServer(app)

// Listen on provided port, on all network interfaces.
server.listen(port)
server.on('error', onError)
console.log('Mock server started on port ' + port + '!')
console.log('简历管理系统 API 已注册!')

// Event listener for HTTP server "error" event.
function onError(error: any) {
  if (error.syscall !== 'listen') {
    throw error
  }
  const bind = typeof port === 'string' ? 'Pipe ' + port : 'Port ' + port
  // handle specific listen errors with friendly messages
  switch (error.code) {
    case 'EACCES':
      console.error('Express ERROR (app) : %s requires elevated privileges', bind)
      process.exit(1)
    case 'EADDRINUSE':
      console.error('Express ERROR (app) : %s is already in use', bind)
      process.exit(1)
    default:
      throw error
  }
}
