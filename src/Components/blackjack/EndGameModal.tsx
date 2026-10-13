import { useRef,useEffect } from "react";


type EndGameModalProps = {
    message: string;
    open: boolean;
    resetGame: () => void;
    goHome: () => void;
  };

export default function EndGameModal({message, open, resetGame, goHome} :EndGameModalProps){

    const dialogRef = useRef<HTMLDialogElement>(null);

    useEffect(() => {
        if (open) {
          dialogRef.current?.showModal();
        } else {
          dialogRef.current?.close();
        }
      }, [open]);
    return (
        <>
        <dialog className="endGame" ref={dialogRef}>
          <h2>{message}</h2>
          <div className="button-holder">
          <button onClick={resetGame}>Play Again</button>
          <button onClick={goHome}>Go Home</button>
          </div>
        </dialog>

        </>

    )
}