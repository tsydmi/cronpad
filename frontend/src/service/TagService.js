import axios from "axios";

export default {
    create(tag) {
        return axios.post('/admin/tags', this.convertToBody(tag))
    },
    findAll() {
        return axios.get('/tags')
    },
    update(tag) {
        return axios.put(`/admin/tags/${tag.id}`, this.convertToBody(tag))
    },
    delete(tag) {
        return axios.delete(`/admin/tags/${tag.id}`)
    },
    convertToBody(tag) {
        return {
            name: tag.name,
            description: tag.description,
            color: tag.color.hex ? tag.color.hex : tag.color,
        }
    }
}