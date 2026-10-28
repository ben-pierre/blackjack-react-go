import { useEffect, useState } from "react";
import { Link } from "react-router";
import {
  getCurrentUser,
  signInWithRedirect,
  signOut,
  fetchAuthSession,
} from "aws-amplify/auth";
import { Hub } from "aws-amplify/utils";

function AuthBar() {
  const [user, setUser] = useState<{ username: string } | null>(null);
  const [email, setEmail] = useState<string | undefined>();

  useEffect(() => {
    async function refreshUser() {
      try {
        const currentUser = await getCurrentUser();
        setUser({ username: currentUser.username });

        try {
          const session = await fetchAuthSession();
          const email = session.tokens?.idToken?.payload.email;
          setEmail(typeof email === "string" ? email : undefined);
        } catch {
          setEmail("Could not fetch");
        }
      } catch {
        setUser(null);
        setEmail(undefined);
      }
    }

    void refreshUser();

    const cancel = Hub.listen("auth", ({ payload }) => {
      if (payload.event === "signedIn") {
        void refreshUser();
      } else if (payload.event === "signedOut") {
        setUser(null);
        setEmail(undefined);
      }
    });

    return cancel;
  }, []);

  return (
    <header className="auth-bar">
      <nav className="auth-nav">
        <Link to="/">Home</Link>
        <Link to="/blackjack">Blackjack</Link>
        <Link to="/Account">Account</Link>
      </nav>

      <div className="auth-user">
        {user ? (
          <>
            Signed in as <Link to="/account"> {email ?? user.username}</Link>
          </>
        ) : (
          "Playing as guest"
        )}

        {user ? (
          <button onClick={() => void signOut()} className="actionButton">
            Log out
          </button>
        ) : (
          <button
            onClick={() => void signInWithRedirect()}
            className="actionButton"
          >
            Log in
          </button>
        )}
      </div>
    </header>
  );
}

export default AuthBar;