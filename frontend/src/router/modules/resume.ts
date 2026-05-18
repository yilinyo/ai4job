import { RouteConfig } from 'vue-router'
import Layout from '@/layout/index.vue'

const resumeRoutes: RouteConfig = {
  path: '/resume',
  component: Layout,
  redirect: '/resume/repos',
  name: 'Resume',
  meta: {
    title: 'resumeManagement',
    icon: 'documentation'
  },
  children: [
    {
      path: 'repos',
      component: () => import(/* webpackChunkName: "resume-repo-list" */ '@/views/resume-repo/list.vue'),
      name: 'ResumeRepoList',
      meta: { title: 'resumeRepository', icon: 'documentation' }
    },
    {
      path: 'repo/:id',
      component: () => import(/* webpackChunkName: "resume-repo-detail" */ '@/views/resume-repo/detail.vue'),
      name: 'ResumeRepoDetail',
      meta: { title: 'resumeRepositoryDetail', hidden: true }
    },
    {
      path: 'version/create',
      component: () => import(/* webpackChunkName: "resume-version-editor" */ '@/views/resume-version/editor.vue'),
      name: 'ResumeVersionCreate',
      meta: { title: 'createResumeVersion', hidden: true }
    },
    {
      path: 'version/:id',
      component: () => import(/* webpackChunkName: "resume-version-view" */ '@/views/resume-version/view.vue'),
      name: 'ResumeVersionView',
      meta: { title: 'viewResumeVersion', hidden: true }
    },
    {
      path: 'version/:id/edit',
      component: () => import(/* webpackChunkName: "resume-version-editor" */ '@/views/resume-version/editor.vue'),
      name: 'ResumeVersionEdit',
      meta: { title: 'editResumeVersion', hidden: true }
    },
    {
      path: 'application/create',
      component: () => import(/* webpackChunkName: "resume-application-editor" */ '@/views/resume-application/editor.vue'),
      name: 'ResumeApplicationCreate',
      meta: { title: 'createApplicationRecord', hidden: true, activeMenu: '/resume/repos' }
    },
    {
      path: 'application/:id',
      component: () => import(/* webpackChunkName: "resume-application-view" */ '@/views/resume-application/view.vue'),
      name: 'ResumeApplicationView',
      meta: { title: 'applicationRecordDetail', hidden: true, activeMenu: '/resume/repos' }
    },
    {
      path: 'application/:id/edit',
      component: () => import(/* webpackChunkName: "resume-application-editor" */ '@/views/resume-application/editor.vue'),
      name: 'ResumeApplicationEdit',
      meta: { title: 'editApplicationRecord', hidden: true, activeMenu: '/resume/repos' }
    }
  ]
}

export const jobPostingRoutes: RouteConfig = {
  path: '/job-postings',
  component: Layout,
  redirect: '/job-postings/list',
  name: 'JobPostings',
  meta: {
    title: 'jobRepository',
    icon: 'shopping'
  },
  children: [
    {
      path: 'list',
      component: () => import(/* webpackChunkName: "job-posting-list" */ '@/views/job-posting/list.vue'),
      name: 'JobPostingList',
      meta: { title: 'jobRepository', icon: 'shopping' }
    },
    {
      path: 'detail/:id',
      component: () => import(/* webpackChunkName: "job-posting-detail" */ '@/views/job-posting/detail.vue'),
      name: 'JobPostingDetail',
      meta: { title: 'jobRepositoryDetail', hidden: true, activeMenu: '/job-postings/list' }
    }
  ]
}

export const interviewQuestionBankRoutes: RouteConfig = {
  path: '/interview-question-bank',
  component: Layout,
  redirect: '/interview-question-bank/list',
  name: 'InterviewQuestionBank',
  meta: {
    title: 'interviewQuestionBank',
    icon: 'clipboard',
    alwaysShow: true
  },
  children: [
    {
      path: 'list',
      component: () => import(/* webpackChunkName: "interview-question-bank" */ '@/views/interview-question-bank/index.vue'),
      name: 'InterviewQuestionBankList',
      meta: { title: 'interviewQuestionLibrary' }
    },
    {
      path: 'flashcard',
      component: () => import(/* webpackChunkName: "interview-flashcard-review" */ '@/views/interview-question-bank/flashcard-review-tab.vue'),
      name: 'InterviewFlashcardReview',
      meta: { title: 'interviewFlashcardReview' }
    },
    {
      path: 'question/:id',
      component: () => import(/* webpackChunkName: "interview-question-detail" */ '@/views/interview-question-bank/detail.vue'),
      name: 'InterviewQuestionDetail',
      meta: { title: 'interviewQuestionDetail', hidden: true, activeMenu: '/interview-question-bank/list' }
    }
  ]
}

export default resumeRoutes
