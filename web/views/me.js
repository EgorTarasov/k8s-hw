import { api, auth } from '/app.js';

const { ref, onMounted } = Vue;
const { useRouter } = VueRouter;

export const Me = {
  template: `
    <article>
      <header><strong>Profile</strong></header>
      <p v-if="loading" aria-busy="true">Loading...</p>
      <dl v-else-if="user">
        <dt>ID</dt><dd><code>{{ user.id }}</code></dd>
        <dt>Email</dt><dd>{{ user.email }}</dd>
        <dt>Name</dt><dd>{{ user.name }}</dd>
      </dl>
      <p v-if="error" class="error">{{ error }}</p>
    </article>
  `,
  setup() {
    const router = useRouter();
    const user = ref(null);
    const error = ref('');
    const loading = ref(true);

    onMounted(async () => {
      try {
        user.value = await api('/api/me');
      } catch (e) {
        error.value = e.message;
        if (e.status === 401) {
          auth.setToken('');
          router.push('/login');
        }
      } finally {
        loading.value = false;
      }
    });
    return { user, error, loading };
  },
};
