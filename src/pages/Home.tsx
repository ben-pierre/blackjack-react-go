import discordLogo from "./../assets/discord.svg";
import "./Home.css";
import { useState, useEffect } from "react";
import { Link } from "react-router";

export default function Home() {
  const reasons = [
    "The government is shut down",
    "Too many cheating competitors",
    "The economy is down",
    "Prime day is next month",
    "Amazon is suppressing you",
  ];

  const [reasonStates, setReasonStates] = useState(
    reasons.map(() => ({
      visible: false,
      crossed: false,
    }))
  );

  const [showResult, setShowResult] = useState(false);
  const [showButtons, setShowButtons] = useState(false);

  function playAnimation() {
    const step = 960;
    reasons.forEach((_, i) => {
      setTimeout(() => {
        setReasonStates(previous=>{
            const next = [...previous];
            next[i] = {...next[i], visible:true}
            return next;
        })
      }, i * step);
      setTimeout(() => {
        setReasonStates(prev => {
          const next = [...prev];
          next[i] = { ...next[i], crossed: true };
          return next;
        });
      }, i * step + step * 0.75);
      setTimeout(()=>{
        console.log
        setShowResult(true);
      }, (reasons.length) * step)
      setTimeout(()=>{
        console.log
        setShowButtons(true);
      }, (reasons.length+1) * step)


    });
  }

  useEffect(() => {
    playAnimation();
  }, []);

  return (
    <div className="v-center">
      <h1>Why Are Sales Down?</h1>
      <ul id="reasons">
        {reasons.map((reason, i) => (
          <li
            key={i}
            className={`${reasonStates[i].visible ? "visible" : "hidden"} ${
              reasonStates[i].crossed ? "crossed" : ""
            }`}
          >
            {reason}
          </li>
        ))}
      </ul>

      <div aria-live="polite" className={showResult?"result visible":"result hidden"}>
        It's you
      </div>

      <div aria-live="polite" className={showButtons?"button-holder visible":"button-holder hidden"} >

      <a
        href="https://discord.gg/KMW323z2DY"
        target="_blank"
        className="button blue"
      >
        <img src={discordLogo} alt="Discord logo" />
        Join the largest online Amazon FBA community
      </a>
      <Link to="/blackjack" className="button red"> Play Blackjack!</Link>
      </div>
      
    </div>
  );
}
