import Link from 'next/link';

export default function HomePage() {
  return (
    <main className="flex h-screen flex-col justify-center text-center">
      <h1 className="mb-4 text-2xl font-bold">Hello World</h1>
      <div className='flex items-center justify-center space-x-10'>
        <div>
          <h1>
            <Link href={"/docs/golang"}>
              Golang Book
            </Link>
          </h1>
        </div>
        <div>
          <h1>
            <Link href={"/docs/office"}>
              Office Book
            </Link>
          </h1>
        </div>
      </div>
    </main>
  );
}
