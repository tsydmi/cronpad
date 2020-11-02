<template>
  <v-card>
    <v-card-title>
      <span class="headline">Configure tags</span>
    </v-card-title>
    <v-card-text>
      <v-container>
        <v-row
            v-for="tag in tags"
            v-bind:key="tag.id"
            align="center"
            justify="start"
            class="ma-1"
        >
          <span>
            {{ tag.name }}
          </span>
          <v-spacer></v-spacer>
          <v-btn
              class="error--text"
              text
              @click="deleteTag(tag)"
          >
            REMOVE
          </v-btn>
        </v-row>
      </v-container>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
          color="blue darken-1"
          text
          @click="closeDialog"
      >
        Close
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import axios from 'axios'

export default {
  name: "ConfigureTags",
  props: {
    dialog: Boolean,
    tags: Array
  },
  data: () => ({}),
  methods: {
    deleteTag(tag) {
      console.log('deleting')
      axios.delete(`/tags/${tag.id}`)
          .then((response) => {
            if (response.status >= 200 && response.status < 300)
              this.$emit('refreshTags', null)
          })
    },
    closeDialog() {
      this.$emit('closeDialog', null)
    }
  }
}
</script>

<style scoped>

</style>