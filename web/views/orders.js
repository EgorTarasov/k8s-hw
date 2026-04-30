import { api } from '/app.js';

const { ref } = Vue;
const { useRouter } = VueRouter;

export const Orders = {
  template: `
    <article>
      <header><strong>New order</strong></header>
      <form @submit.prevent="submit">
        <div v-for="(it, idx) in items" :key="idx" class="row">
          <label>SKU
            <input v-model="it.sku" required />
          </label>
          <label>Qty
            <input v-model.number="it.qty" type="number" min="1" required />
          </label>
          <button type="button" class="secondary" @click="remove(idx)" :disabled="items.length === 1">×</button>
        </div>
        <button type="button" class="secondary" @click="add">+ Add item</button>
        <button :aria-busy="loading" type="submit">Create order</button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </article>
  `,
  setup() {
    const router = useRouter();
    const items = ref([{ sku: '', qty: 1 }]);
    const error = ref('');
    const loading = ref(false);

    const add = () => items.value.push({ sku: '', qty: 1 });
    const remove = (idx) => items.value.splice(idx, 1);

    const submit = async () => {
      loading.value = true;
      error.value = '';
      try {
        const out = await api('/api/orders', { method: 'POST', body: { items: items.value } });
        router.push(`/orders/${out.id}`);
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };
    return { items, error, loading, add, remove, submit };
  },
};
