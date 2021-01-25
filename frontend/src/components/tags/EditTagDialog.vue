<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="400px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">Update a new Tag</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-form ref="form" v-model="valid">
            <v-row>
              <v-text-field
                  label="Name"
                  v-model="tag.name"
                  @keydown.enter="saveTag"
                  :rules="rules.name"
              ></v-text-field>
            </v-row>
            <v-row>
              <v-text-field
                  label="Description"
                  v-model="tag.description"
                  @keydown.enter="saveTag"
                  required
              ></v-text-field>
            </v-row>
            <v-row>
              <tag-color-picker v-model="tag.color"/>
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
import TagService from "@/service/TagService"
import TagColorPicker from "@/components/tags/TagColorPicker"

export default {
  name: "CreateTagDialog",
  components: {
    TagColorPicker,
  },

  props: {
    value: {
      type: Boolean,
      required: false,
    },
    tag: {
      type: Object,
      required: true,
    },
  },
  data: () => ({
    valid: true,
    rules: {
      name: [
        v => !!v || 'Name is required',
      ],
    },
  }),
  methods: {
    saveTag() {
      if (this.$refs.form.validate()) {
        TagService.update(this.tag)
            .then(() => {
              this.$emit('refreshTags', null)
              this.emitChange(false)
            })
      }
    },
    emitChange(value) {
      this.$emit('input', value)
    }
  },
}
</script>