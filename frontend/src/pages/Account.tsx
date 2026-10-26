import {
  getCurrentUser,
  type AuthUser,
  signInWithRedirect,
  signOut,
  fetchAuthSession,
} from "aws-amplify/auth";
import { useEffect, useState } from "react";

export default function Account() {
  const [currentUser, setCurrentUser] = useState<AuthUser | null>(null);
  const [email, setEmail] = useState<string | undefined>();
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    async function loadUser() {
      try {
        const user = await getCurrentUser();
        setCurrentUser(user);
      } catch (error) {
        console.error("Failed to get current user:", error);
        setCurrentUser(null);
      }
      try {
        const session = await fetchAuthSession();
        const email = session.tokens?.idToken?.payload.email;
        setEmail(typeof email === "string" ? email : undefined);
      } catch (error) {
        console.error("Failed to get current user attributes:", error);
        setEmail("Could not fetch");
      }
      setLoading(false);
    }

    loadUser();
  }, []);

  return (
    <>
      {loading ? (
        <p><div className="spinner" aria-label="Loading account" />Loading account...</p>

      ) : currentUser == null ? (
        <div>
          <p>Not signed in.</p>
          <button onClick={() => signInWithRedirect()}>Log in</button>
        </div>
      ) : (
        <div>
          <p>Hello '{email ?? "Not available"}'!</p>
          <p>Your accountId is {currentUser.userId}!</p>
          <button onClick={() => signOut()}>Sign out</button>
        </div>
      )}
    </>
  );
}
