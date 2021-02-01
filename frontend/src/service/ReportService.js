import axios from "axios";

export default {
    findAll(query) {
        return axios.post('/admin/user-reports', query)
    },
    getProjectReport(projectID) {
        return axios.get(`/manager/project-reports/${projectID}`)
    },
}