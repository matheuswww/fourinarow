"use client"
import { SyntheticEvent, useState } from "react"
import styles from "./form.module.css"
import Link from "next/link";
import { useRouter } from "next/navigation";
import { basePath } from "@/api/path";


interface Props {
  mode: 'signin' | 'signup';
}

export default function Form({ mode }: Props) {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter()

  async function handleSignup(event: SyntheticEvent) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(basePath+'/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_name: username,
          password: password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        switch (response.status) {
          case 400:
            if (data.message === 'invalid fields') {
              throw new Error('Invalid fields. Please check the data provided.');
            } else if (data.message === 'user already exists') {
              throw new Error('User already exists.');
            } else if (data.message === 'server error') {
              throw new Error('Server error while generating token.');
            }
            break;
          case 500:
            throw new Error('Internal server error.');
          default:
            throw new Error(`Unknown error: ${data.message || 'Please try again.'}`);
        }
      }

      console.log('Signup successful!');
      console.log('Token:', data.token);
      console.log('User ID:', data.user_id);
      window.localStorage.setItem("token", data.token)
      window.localStorage.setItem("user_id", data.user_id)
      router.push("/")
    } catch (error: any) {
      console.error('Error during signup:', error.message);
      setError(error.message);
      return null;
    } finally {
      setLoading(false);
    }
  }

  async function handleSignin(event: SyntheticEvent) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(basePath+'/signin', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_name: username,
          password: password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        switch (response.status) {
          case 400:
            if (data.message === 'invalid fields') {
              throw new Error('Invalid fields. Please check the data provided.');
            } else if (data.message === 'user not found') {
              throw new Error('User not found.');
            } else if (data.message === 'invalid password') {
              throw new Error('Invalid password.');
            } else if (data.message === 'server error') {
              throw new Error('Server error while generating token.');
            }
            break;
          case 500:
            throw new Error('Internal server error.');
          default:
            throw new Error(`Unknown error: ${data.message || 'Please try again.'}`);
        }
      }

      console.log('Login successful!');
      console.log('Token:', data.token);
      console.log('User ID:', data.user_id);
      window.localStorage.setItem("token", data.token)
      window.localStorage.setItem("user_id", data.user_id)
      router.push("/")
    } catch (error: any) {
      console.error('Error during signin:', error.message);
      setError(error.message);
      return null;
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className={styles.section}>
      <form
        className={styles.form}
        onSubmit={mode === 'signin' ? handleSignin : handleSignup}
      >
        <label htmlFor="user_name">User Name</label>
        <input
          type="text"
          id="user_name"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          disabled={loading}
        />

        <label htmlFor="password">Password</label>
        <input
          type="password"
          id="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          disabled={loading}
        />
        { mode == "signup" ? <Link className={styles.link} href={"/signin"}>Já possui uma conta?</Link> : <Link className={styles.link} href={"/signup"}>Não possui uma conta?</Link> }
        {error && <p className={styles.error}>{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? 'Loading...' : mode === 'signin' ? 'Sign In' : 'Sign Up'}
        </button>
      </form>
    </section>
  );
}