import AuthButton from "~/app/auth-button";
import NextAuthProvider from "~/app/next-auth-provider";

export default async function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center p-24 gap-2">
      <NextAuthProvider>
        <AuthButton />
      </NextAuthProvider>
    </main>
  );
}
