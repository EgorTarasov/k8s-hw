import { api } from '/app.js';

const { ref } = Vue;
const { useRouter } = VueRouter;

export const Register = {
  template: `
    <article>
      <header><strong>Register</strong></header>
      <form @submit.prevent="submit">
        <label>Name
          <input v-model="name" required />
        </label>
        <label>Email
          <input v-model="email" type="email" autocomplete="email" required />
        </label>
        <label>Password (min 6 chars)
          <input v-model="password" type="password" minlength="6" autocomplete="new-password" required />
        </label>
        <button :aria-busy="loading" type="submit">Create account</button>
        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="success" class="muted">Registered! Redirecting to login...</p>
      </form>
      <footer><router-link to="/login">Already registered?</router-link></footer>
    </article>
  `,
  setup() {
    const router = useRouter();
    const name = ref('');
    const email = ref('');
    const password = ref('');
    const error = ref('');
    const success = ref(false);
    const loading = ref(false);

    const submit = async () => {
      loading.value = true;
      error.value = '';
      try {
        await api('/api/register', { method: 'POST', body: { name: name.value, email: email.value, password: password.value } });
        success.value = true;
        setTimeout(() => router.push('/login'), 600);
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };
    return { name, email, password, error, success, loading, submit };
  },
};
