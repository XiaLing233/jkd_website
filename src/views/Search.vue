<template>
  <div>
    <h1>课程检索</h1>
    <div v-for="(condition, index) in conditions" :key="index" class="search-condition">
      <!-- Show the connector dropdown only if it's not the first condition -->
      <div v-if="index !== 0">
        <span>Connector:</span>
        <select v-model="condition.connector">
          <option value="and">AND</option>
          <option value="or">OR</option>
          <option value="not">NOT</option>
        </select>
      </div>

      <select v-model="condition.selectedItem">
        <option value="" disabled>Select Item</option>
        <option v-for="item in availableItemsFor(index)" :key="item" :value="item">{{ item }}</option>
      </select>
      <input
        type="text"
        v-model="condition.searchWord"
        placeholder="Enter search word"
      />

      <button @click="removeCondition(index)">Remove</button>
    </div>
    <button @click="addCondition" :disabled="!canAddCondition">Add Condition</button>
    <button @click="submitSearch">Search</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      conditions: [{ selectedItem: '', searchWord: '', connector: '' }],
      allItems: ['Item 1', 'Item 2', 'Item 3', 'Item 4'], // All available items
    };
  },
  computed: {
    canAddCondition() {
      // Check if the number of conditions is less than all items
      return this.conditions.length < this.allItems.length;
    },
  },
  methods: {
    addCondition() {
      this.conditions.push({ selectedItem: '', searchWord: '', connector: '' });
    },
    removeCondition(index) {
      this.conditions.splice(index, 1);
    },
    // updateAvailableItems(index) {
    //   // No specific updates needed here as we're managing items separately
    // },
    availableItemsFor(index) {
      // Get all the selected items except the current one
      const selectedItems = this.conditions
        .filter((cond, i) => i !== index)
        .map(cond => cond.selectedItem);
      // Return all items except those already selected by other conditions
      return this.allItems.filter(item => !selectedItems.includes(item));
    },
    submitSearch() {
      const requestData = this.conditions.filter(cond => cond.selectedItem); // Filter out empty conditions
      // Sending the HTTP request to the Flask backend
      fetch('http://localhost:5000/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      })
        .then(response => response.json())
        .then(data => {
          console.log('Search results:', data);
        })
        .catch(error => {
          console.error('Error:', error);
        });
    },
  },
};
</script>

<style scoped>
.search-condition {
  margin-bottom: 10px;
}
input[type="text"] {
  margin-left: 10px;
  margin-right: 10px;
}
</style>
