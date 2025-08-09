import {Layout} from '../components/Layout';
import {Container} from '../components/ui/Container';
import {FadeIn} from '../components/ui/FadeIn';
import {LinkButton} from '../components/ui/Button';

export default function LandingPage() {
    return (
        <Layout>
            <Container className="mt-24 sm:mt-32 md:mt-56">
                <FadeIn className="max-w-3xl text-center mx-auto">
                    <h1 className="font-display text-5xl font-medium tracking-tight text-balance text-neutral-950 sm:text-7xl">
                        StrikePad
                    </h1>
                    <p className="mt-6 text-xl text-neutral-600">
                        Coming Soon
                    </p>
                    <div className="mt-8 flex flex-col gap-4 sm:flex-row justify-center">
                        <LinkButton href="/signup" size="lg" className="px-8 py-4">
                            サインアップ
                        </LinkButton>
                        <LinkButton href="/login" variant="outline" size="lg" className="px-8 py-4">
                            ログイン
                        </LinkButton>
                    </div>
                </FadeIn>
            </Container>
        </Layout>
    );
}