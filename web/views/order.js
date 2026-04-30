import { api } from '/app.js';

const { ref, computed, onMounted, onUnmounted } = Vue;

const TERMINAL = new Set(['processed', 'failed']);

export const Order = {
  props: ['id'],
  template: `
    <article>
      <header>
        <strong>Order</strong>
        <code>{{ id }}</code>
      </header>
      <p v-if="loading && !order" aria-busy="true">Loading...</p>
      <template v-else-if="order">
        <p>
          Status:
          <mark :data-tooltip="order.status">{{ order.status }}</mark>
          <span v-if="!terminal" class="muted"> (refreshing...)</span>
        </p>
        <table>
          <thead><tr><th>SKU</th><th>Qty</th></tr></thead>
          <tbody>
            <tr v-for="(it, idx) in order.items" :key="idx"><td>{{ it.sku }}</td><td>{{ it.qty }}</td></tr>
          </tbody>
        </table>
      </template>
      <p v-if="error" class="error">{{ error }}</p>
      <footer><router-link to="/orders">Create another</router-link></footer>
    </article>
  `,
  setup(props) {
    const order = ref(null);
    const error = ref('');
    const loading = ref(true);
    let timer = null;

    const terminal = computed(() => order.value && TERMINAL.has(order.value.status));

    const refresh = async () => {
      try {
        order.value = await api(`/api/orders/${props.id}`);
        error.value = '';
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
      if (terminal.value && timer) {
        clearInterval(timer);
        timer = null;
      }
    };

    onMounted(() => {
      refresh();
      timer = setInterval(refresh, 1000);
    });
    onUnmounted(() => timer && clearInterval(timer));

    return { order, error, loading, terminal };
  },
};
