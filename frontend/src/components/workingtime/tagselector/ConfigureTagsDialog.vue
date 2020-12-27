<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="400px"
  >
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
            @click="emitChange(false)"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import TagService from "@/service/TagService";

export default {
  name: "ConfigureTagsDialog",
  props: {
    value: {
      type: Boolean,
      required: false,
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
    emitChange(value) {
      this.$emit('input', value)
    },
  }
}
</script>