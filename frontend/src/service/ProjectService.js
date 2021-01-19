import axios from "axios";

export default {
    findCurrentUserProjects() {
        return axios.get('/projects')
    },
    create(project) {
        return axios.post('/admin/projects', project)
    },
    getUsers(projectID) {
        return axios.get(`/admin/projects/${projectID}/users`)
    },
    search(search) {
        return axios.post(`/admin/projects/search`, search)
    },
    findAll() {
        return axios.post(`/admin/projects/search`, {})
    },
    getProjectDetails(project) {
        return axios.get(`/admin/projects/${project.id}`)
    },
    update(project) {
        return axios.put(`/admin/projects/${project.id}`, project)
    },
    delete(projectID) {
        return axios.delete(`/admin/projects/${projectID}`)
    },
}