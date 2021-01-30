import axios from "axios"
import cloneDeep from 'clone-deep'

export default {
    create(tag) {
        if (tag.basic) {
            return axios.post('/admin/base-tags', this.convertToBody(tag))
        } else {
            return axios.post('/manager/tags', this.convertToBody(tag))
        }
    },
    findAll() {
        return axios.get('/tags')
    },
    update(tag) {
        if (tag.basic) {
            return axios.put(`/admin/base-tags/${tag.id}`, this.convertToBody(tag))
        } else {
            return axios.put(`/manager/tags/${tag.id}`, this.convertToBody(tag))
        }
    },
    delete(tag) {
        if (tag.basic) {
            return axios.delete(`/admin/base-tags/${tag.id}`)
        } else {
            return axios.delete(`/manager/tags/${tag.id}`)
        }
    },

    convertToBody(tag) {
        let body = cloneDeep(tag)
        body.color = tag.color.hex ? tag.color.hex : tag.color
        return body
    }
}