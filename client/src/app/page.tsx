import styles from "./page.module.css";
import Player from "@/components/Player";

export default function Home() {
  return (
    <div className={styles.page}>
      <Player />
    </div>
  );
}
