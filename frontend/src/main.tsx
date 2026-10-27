import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { Amplify } from 'aws-amplify';
import "aws-amplify/auth/enable-oauth-listener";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)



Amplify.configure({
  Auth: {
    Cognito: {
      userPoolId: import.meta.env.VITE_USER_POOL_ID,
      userPoolClientId: import.meta.env.VITE_APP_POOL_ID,
      loginWith: {
        oauth: {
          domain: 'auth.whyaresalesdown.com',
          scopes: ['openid', 'email', 'phone'],
          redirectSignIn: [
            'http://localhost:5173/blackjack',
            'https://whyaresalesdown.com/blackjack',
          ],
          redirectSignOut: [
            'http://localhost:5173/blackjack',
            'https://whyaresalesdown.com/blackjack',
          ],
          responseType: 'code',
        },
      },
    },
  },
});