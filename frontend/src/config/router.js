import Vue from 'vue'
import VueRouter from 'vue-router'

import Home from "@/pages/Home";
import WorkingTime from "@/pages/WorkingTime";
import Settings from "@/pages/Settings";
import About from "@/pages/About";
import Page404 from "@/pages/Page404";
import Projects from "@/pages/Projects";
import Reports from "@/pages/UserReports";
import Tags from "@/pages/Tags";

const router = new VueRouter({
    mode: 'history',
    routes: [
        {path: '/', component: Home},
        {path: '/working-time', component: WorkingTime},
        {path: '/projects', component: Projects},
        {path: '/user-reports', component: Reports},
        {path: '/tags', component: Tags},
        {path: '/settings', component: Settings},
        {path: '/about', component: About},
        {path: "*", component: Page404}
    ]
});

Vue.use(VueRouter)

export default router