import AuthButton from "~/components/auth-button";
import NextAuthProvider from "~/auth/next-auth-provider";

export default async function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center p-24 gap-2">
      <NextAuthProvider>
        <AuthButton />
      </NextAuthProvider>
    </main>
  );
}
