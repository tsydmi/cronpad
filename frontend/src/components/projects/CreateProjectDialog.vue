<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="600px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">Create a new Project</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-form
              ref="form"
              lazy-validation
          >
            <v-row>
              <v-text-field
                  label="Name"
                  v-model="project.name"
                  :rules="rules.name"
                  required
              ></v-text-field>
            </v-row>

            <v-row>
              <v-text-field
                  label="Description"
                  v-model="project.description"
                  :rules="rules.description"
                  required
              ></v-text-field>
            </v-row>

            <v-row class="pt-4 pb-1">
              <user-select v-model="project.users" :users="users"/>
            </v-row>

          </v-form>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-btn
            color="primary"
            text
            @click="save"
        >
          Save
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="primary"
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
import ProjectService from "@/service/ProjectService"
import UserSelect from "@/components/UserSelect"

export default {
  name: "CreateProjectDialog",
  components: {
    UserSelect
  },

  props: {
    value: {
      type: Boolean,
      required: false,
    },
    users: {
      type: Array,
      required: true,
    }
  },
  data: () => ({
    project: {},
    rules: {
      name: [
        v => !!v || 'Name is required',
      ],
      description: [
        v => !!v || 'Description is required',
      ],
    },
  }),
  methods: {
    save() {
      if (this.$refs.form.validate()) {
        ProjectService.create(this.project)
            .then(() => {
              this.$emit('refreshProjects', null)
              this.emitChange(false)
            })
      }
    },
    emitChange(value) {
      this.$emit('input', value)
    },
  },
  watch: {
    'value': function () {
      if (this.$refs.form) {
        this.$refs.form.reset()
      }
    }
  }
}
</script>