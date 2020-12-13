import axios from "axios";

export default {
    create(tag) {
        let body = {
            name: tag.name,
            parent: tag.parent,
            color: tag.color.hex,
        }

        return axios.post('/tags', body)
    },
    findAll() {
        return axios.get('/tags')
    },
    delete(tag) {
        return axios.delete(`/tags/${tag.id}`)
    },
}