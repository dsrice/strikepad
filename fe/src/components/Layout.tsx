import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { clsx } from 'clsx';
import { Container } from './ui/Container';
import { LinkButton } from './ui/Button';

interface LayoutProps {
  children: React.ReactNode;
}

function MenuIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true" className={className}>
      <path d="M2 6h20v2H2zM2 16h20v2H2z" />
    </svg>
  );
}

function XIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true" className={className}>
      <path d="m5.636 4.223 14.142 14.142-1.414 1.414L4.222 5.637z" />
      <path d="M4.222 18.363 18.364 4.22l1.414 1.414L5.636 19.777z" />
    </svg>
  );
}

function Logo({ className }: { className?: string }) {
  return (
    <div className={clsx('font-bold text-xl', className)}>
      StrikePad
    </div>
  );
}

function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  return (
    <header className="relative z-50">
      {/* Main Header */}
      <div className="absolute top-2 right-0 left-0 z-40 pt-14">
        <Container>
          <div className="flex items-center justify-between">
            <Link
              to="/"
              aria-label="Home"
              className="text-neutral-950 hover:text-blue-600 transition-colors"
            >
              <Logo />
            </Link>
            <div className="flex items-center gap-x-8">
              <LinkButton href="/login" variant="outline" size="sm">
                ログイン
              </LinkButton>
              <LinkButton href="/signup" variant="primary" size="sm">
                サインアップ
              </LinkButton>
              <button
                type="button"
                onClick={() => setIsMenuOpen(!isMenuOpen)}
                aria-expanded={isMenuOpen}
                className={clsx(
                  'group -m-2.5 rounded-full p-2.5 transition hover:bg-neutral-950/10'
                )}
                aria-label="Toggle navigation"
              >
                {isMenuOpen ? (
                  <XIcon className="h-6 w-6 fill-neutral-950 group-hover:fill-neutral-700" />
                ) : (
                  <MenuIcon className="h-6 w-6 fill-neutral-950 group-hover:fill-neutral-700" />
                )}
              </button>
            </div>
          </div>
        </Container>
      </div>

      {/* Navigation Menu */}
      {isMenuOpen && (
        <div className="fixed inset-0 z-50 bg-neutral-950 pt-2">
          <div className="bg-neutral-800">
            <div className="bg-neutral-950 pt-14 pb-16">
              <Container>
                <div className="flex items-center justify-between">
                  <Link
                    to="/"
                    aria-label="Home"
                    className="text-white hover:text-neutral-200 transition-colors"
                  >
                    <Logo className="text-white" />
                  </Link>
                  <button
                    type="button"
                    onClick={() => setIsMenuOpen(false)}
                    className="group -m-2.5 rounded-full p-2.5 transition hover:bg-white/10"
                    aria-label="Close navigation"
                  >
                    <XIcon className="h-6 w-6 fill-white group-hover:fill-neutral-200" />
                  </button>
                </div>
              </Container>
            </div>
          </div>
          <nav className="mt-px font-display text-5xl font-medium tracking-tight text-white">
            <div className="even:mt-px sm:bg-neutral-950">
              <Container>
                <div className="grid grid-cols-1 sm:grid-cols-2">
                  <Link
                    to="/features"
                    className="group relative isolate -mx-6 bg-neutral-950 px-6 py-10 even:mt-px sm:mx-0 sm:px-0 sm:py-16 sm:odd:pr-16 sm:even:mt-0 sm:even:border-l sm:even:border-neutral-800 sm:even:pl-16"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    機能
                    <span className="absolute inset-y-0 -z-10 w-screen bg-neutral-900 opacity-0 transition group-odd:right-0 group-even:left-0 group-hover:opacity-100" />
                  </Link>
                  <Link
                    to="/pricing"
                    className="group relative isolate -mx-6 bg-neutral-950 px-6 py-10 even:mt-px sm:mx-0 sm:px-0 sm:py-16 sm:odd:pr-16 sm:even:mt-0 sm:even:border-l sm:even:border-neutral-800 sm:even:pl-16"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    料金
                    <span className="absolute inset-y-0 -z-10 w-screen bg-neutral-900 opacity-0 transition group-odd:right-0 group-even:left-0 group-hover:opacity-100" />
                  </Link>
                </div>
              </Container>
            </div>
          </nav>
        </div>
      )}
    </header>
  );
}

function Footer() {
  return (
    <footer className="mt-24 w-full sm:mt-32 lg:mt-40">
      <Container>
        <div className="border-t border-neutral-200 pt-16 pb-8">
          <div className="flex flex-col items-center justify-between gap-y-12 lg:flex-row lg:gap-y-0">
            <Link to="/" aria-label="Home">
              <Logo />
            </Link>
            <nav className="flex flex-wrap gap-x-8 gap-y-2 justify-center lg:justify-end">
              <Link
                to="/privacy"
                className="text-sm text-neutral-600 hover:text-neutral-900 transition-colors"
              >
                プライバシーポリシー
              </Link>
              <Link
                to="/terms"
                className="text-sm text-neutral-600 hover:text-neutral-900 transition-colors"
              >
                利用規約
              </Link>
              <Link
                to="/contact"
                className="text-sm text-neutral-600 hover:text-neutral-900 transition-colors"
              >
                お問い合わせ
              </Link>
              <Link
                to="/help"
                className="text-sm text-neutral-600 hover:text-neutral-900 transition-colors"
              >
                ヘルプ
              </Link>
            </nav>
          </div>
          <div className="mt-8 flex flex-col items-center justify-between gap-y-4 border-t border-neutral-200 pt-8 lg:flex-row lg:gap-y-0">
            <div className="text-center lg:text-left">
              <p className="text-sm text-neutral-500">
                © {new Date().getFullYear()} StrikePad. All rights reserved.
              </p>
              <p className="text-xs text-neutral-400 mt-1">
                個人開発プロジェクトとして運営されています
              </p>
            </div>
            <div className="text-center lg:text-right">
              <p className="text-sm text-neutral-500">
                Made with ❤️ for better productivity
              </p>
              <p className="text-xs text-neutral-400 mt-1">
                Version 1.0.0
              </p>
            </div>
          </div>
        </div>
      </Container>
    </footer>
  );
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-neutral-950 text-base antialiased">
      <div className="flex min-h-screen flex-col">
        <Header />
        <div 
          className="relative flex flex-auto overflow-hidden bg-white pt-14"
          style={{ borderTopLeftRadius: 40, borderTopRightRadius: 40 }}
        >
          <div className="relative isolate flex w-full flex-col pt-9">
            <main className="w-full flex-1">
              {children}
            </main>
            <Footer />
          </div>
        </div>
      </div>
    </div>
  );
}