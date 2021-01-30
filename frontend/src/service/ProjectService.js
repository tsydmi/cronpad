import axios from "axios";

export default {
    create(project) {
        return axios.post('/admin/projects', project)
    },
    getUsers(projectID) {
        return axios.get(`/admin/projects/${projectID}/users`)
    },
    search(search) {
        return axios.post(`/admin/projects/search`, search)
    },
    findAll(isAdmin = false) {
        if (isAdmin) {
            return axios.post(`/admin/projects/search`, {})
        } else {
            return axios.get('/projects')
        }
    },
    update(project) {
        return axios.put(`/admin/projects/${project.id}`, project)
    },
    delete(projectID) {
        return axios.delete(`/admin/projects/${projectID}`)
    },
}