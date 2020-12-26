import axios from "axios";

export default {
    create(project) {
        return axios.post('/admin/projects', project)
    },
    findCurrentUserProjects() {
        return axios.get('/projects')
    },
    findAll() {
        return axios.get(`/admin/projects`)
    },
    search(search) {
        return axios.post(`/admin/projects/search`, search)
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