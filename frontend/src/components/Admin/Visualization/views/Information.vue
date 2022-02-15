<template>
  <div class="columns is-mobile is-centered" style="overflow-y: scroll; max-height: 500px">
    <div class="column is-half">
      <b-field label="Group">
        <b-select expanded v-model="view.group" placeholder="Group...">
          <option
            v-for="group in allGroups"
            :value="group"
            :key="group.uuid"
          >{{ group.name }}</option>
        </b-select>
      </b-field>
      <b-field label="Index">
        <b-select
          expanded
          v-model="view.selectedIndex"
          placeholder="Index..."
        >
          <option v-for="i in allIndices" :value="i" :key="i.index">{{ i.index }}</option>
        </b-select>
      </b-field>

      <div class="columns" v-if="viewClass == 'line'">
        <div class="column">
          <b-field label="X-Axis">
            <b-select v-model="view.fields[0]"
              placeholder="Select field..."
              expanded>
              <option v-for="field in view.selectedIndex.fields" :value="field.name" :key="field"> {{ field.name }} </option>
            </b-select>
          </b-field>
          <b-field label="Y-Axis">
            <b-select v-model="view.fields[1]"
              placeholder="Select field..."
              expanded>
              <option v-for="field in view.selectedIndex.fields" :value="field.name" :key="field"> {{ field.name }} </option>
            </b-select>
          </b-field>
          </div>
          <div class="column">
          <b-field label="Name">
            <b-input v-model="view.fieldNames[0]" placeholder="Enter a name..."></b-input>
          </b-field>
          <b-field label="Name">
            <b-input v-model="view.fieldNames[1]" placeholder="Enter a name..."></b-input>
          </b-field>
        </div>
      </div>

      <div class="columns" v-if="viewClass == 'pie' || viewClass == 'bar'">
        <div class="column">
          <b-field label="Data">
            <b-select v-model="view.fields[0]"
              placeholder="Select field..."
              expanded>
              <option v-for="field in view.selectedIndex.fields" :value="field.name" :key="field"> {{ field.name }} </option>
            </b-select>
          </b-field>
          </div>
          <div class="column">
          <b-field label="Name">
            <b-input v-model="view.fieldNames[0]" placeholder="Enter a name..."></b-input>
          </b-field>
        </div>
      </div>

      <div v-if="viewClass == 'table'" >
        <div class="columns is-mobile" v-for="(column, index) in tableColumns" :key="column.name">
          <div class="column is-one-fifth">
            <b-field style="padding-top: 30px">
              <b-button type="is-primary" icon-right="delete" @click="tableColumns.splice(index, 1); view.fields.splice(index, 1); view.fieldNames.splice(index, 1)"></b-button>
            </b-field>
          </div>
          <div class="column">
            <b-field label="Field">
              <b-select v-model="view.fields[index]"
                placeholder="Select field..."
                expanded>
                <option v-for="field in view.selectedIndex.fields" :value="field.name" :key="field"> {{ field.name }} </option>
              </b-select>
            </b-field>
            </div>
            <div class="column">
            <b-field label="Name">
              <b-input v-model="view.fieldNames[index]" placeholder="Enter a name..."></b-input>
            </b-field>
          </div>
        </div>
        <button style="margin-top:20px;width:100%;" class="button is-primary" type="button" @click="tableColumns.push({ name: 'field' })"><b-icon icon="plus-circle-outline"></b-icon><span>Add Field</span></button>
      </div>
      <b-field label="Authorized">
        <b-select
          expanded
          v-model="view.authorized"
          placeholder="Authorized Asset..."
        >
          <option v-for="auth in view.group.authorized" :value="auth" :key="auth">{{ auth }}</option>
        </b-select>
      </b-field>
    </div>
  </div>
</template>

<script>
export default {
  props: ['viewClass', 'visualizationToEdit'],
  data () {
    return {
      allGroups: [],
      allAuthorized: [],
      allIndices: [],
      tableColumns: [],
      view: {
        class: "",
        group: {},
        fields: [],
        selectedIndex: { fields: [] },
        fieldNames: [],
        authorized: "",
        name: ""
      }
    }
  },
  mounted () {
    this.fetchIndices()
    this.fetchGroups()
  },
  methods: {
    getView() {
      const viewCopy = JSON.parse(JSON.stringify(this.view))
      viewCopy.class = this.viewClass
      viewCopy.group = this.view.group.uuid
      return viewCopy
    },
    fetchIndices() {
      fetch('/api/fields/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          const array = []
          json.forEach(element => {
            const el = element
            el.fields = element.fields.sort(function (a, b) { return a.name.localeCompare(b.name) })
            array.push(el)
          })
          this.allIndices = array

          if (this.visualizationToEdit.name) {
            this.view = this.visualizationToEdit
            this.view.selectedIndex = this.allIndices.find(e => e.index === this.view.index)
            this.view.fields.forEach(field => {
              this.tableColumns.push({ name: 'field' })
            })
          }
        })
    },
    fetchGroups() {
      fetch('/api/group/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.allGroups.push(json.current)
          if (json.others !== null) {
            json.others.forEach(group => {
              this.allGroups.push(group)
            })
          }
          if (this.visualizationToEdit.class) {
            this.allGroups.forEach(group => {
              if (group.uuid === this.visualizationToEdit.group) {
                this.view.group = group
              }
            })
          }
        })
    }
  }
}
</script>
