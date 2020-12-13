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
import TagService from "@/service/TagService";

export default {
  name: "ConfigureTags",
  props: {
    dialog: {
      type: Boolean,
      required: true,
    },
    tags: {
      type: Array,
      required: true,
    },
  },
  data: () => ({}),
  methods: {
    deleteTag(tag) {
      TagService.delete(tag)
          .then(() => this.$emit('refreshTags', null))
    },
    closeDialog() {
      this.$emit('closeDialog', null)
    }
  }
}
</script>

<style scoped>

</style>