<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="400px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">Create a new Tag</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-form v-model="valid">
            <v-row>
              <v-text-field
                  label="Name"
                  v-model="newTag.name"
                  @keydown.enter="saveTag"
                  required
              ></v-text-field>
            </v-row>
            <v-row>
              <v-select
                  label="Parent"
                  hint="new tag can extend existed one"
                  v-model="newTag.parent"
                  item-text="name"
                  :items="tags"
              ></v-select>
            </v-row>
            <v-row>
              <span class="v-label v-text-field mr-3">Color</span>
              <v-color-picker
                  label="Color"
                  hide-canvas
                  show-swatches
                  :swatches="swatches"
                  swatches-max-height="400px"
                  v-model="newTag.color"
              ></v-color-picker>
            </v-row>
          </v-form>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-btn
            color="blue darken-1"
            text
            @click="saveTag"
        >
          Save
        </v-btn>
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
  name: "CreateTagDialog",
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
  data: () => ({
    newTag: {},
    valid: true,
    swatches: [
      ["#B39DDB", "#9FA8DA", "#90CAF9"],
      ["#81D4FA", "#80DEEA", "#80CBC4"],
      ["#A5D6A7", "#C5E1A5", "#E6EE9C"],
      ["#FFE082", "#FFCC80", "#FFAB91"],
    ]
  }),
  methods: {
    saveTag() {
      TagService.create(this.newTag)
          .then(() => {
            this.$emit('refreshTags', null)
            this.emitChange(false)
          })
    },
    emitChange(value) {
      this.$emit('input', value)
    }
  },
  watch: {
    'value': function () {
      this.newTag = {}
    }
  }
}
</script>