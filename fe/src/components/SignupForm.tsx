import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import GoogleSignSection from './GoogleSignSection';
import {Layout} from './Layout';
import {Container} from './ui/Container';
import {FadeIn} from './ui/FadeIn';

const signupSchema = z.object({
  email: z.string().email('有効なメールアドレスを入力してください'),
  password: z
    .string()
    .min(8, 'パスワードは8文字以上である必要があります')
    .max(50, 'パスワードは50文字以下である必要があります')
    .regex(/[a-z]/, 'パスワードには小文字を含める必要があります')
    .regex(/[A-Z]/, 'パスワードには大文字を含める必要があります')
    .regex(/[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/, 'パスワードには記号を含める必要があります'),
  displayName: z.string().min(1, '表示名を入力してください').max(50, '表示名は50文字以下である必要があります'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'パスワードが一致しません',
  path: ['confirmPassword'],
});

type SignupFormData = z.infer<typeof signupSchema>;

const SignupForm: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const {signup, googleSignup, error} = useAuth();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<SignupFormData>({
    resolver: zodResolver(signupSchema),
  });

  const onSubmit = async (data: SignupFormData) => {
    try {
      setIsLoading(true);
      await signup(data.email, data.password, data.displayName);
      navigate('/dashboard');
    } catch (error: any) {
      setError('root', {
        type: 'manual',
        message: error.message || 'アカウント作成に失敗しました',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleSuccess = async (accessToken: string) => {
    try {
      setIsLoading(true);
      await googleSignup(accessToken);
      navigate('/dashboard');
    } catch (error: any) {
      setError('root', {
        type: 'manual',
        message: error.message || 'Googleサインアップに失敗しました',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleError = (error: string) => {
    // Google認証の準備ができていない場合はエラーメッセージを表示しない
    if (error.includes('準備ができていません') || error.includes('初期化に失敗')) {
      console.warn('Google authentication not available:', error);
      return;
    }

    setError('root', {
      type: 'manual',
      message: error,
    });
  };

  return (
      <Layout>
        <Container className="py-24 sm:py-32">
          <FadeIn className="mx-auto w-full max-w-[600px]">
            <div className="text-center mb-8">
              <h1 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
                StrikePadアカウント作成
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                すでにアカウントをお持ちの方は{' '}
                <button
                    type="button"
                    onClick={() => navigate('/login')}
                    className="font-medium text-blue-600 hover:text-blue-500"
                >
                  ログイン
                </button>
              </p>
            </div>

            {/* Google Sign-up Section */}
            <GoogleSignSection
                onSuccess={handleGoogleSuccess}
                onError={handleGoogleError}
                buttonText="Googleでサインイン"
                disabled={isLoading}
            />
        
        <form className="mt-8 space-y-6" onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                メールアドレス
              </label>
              <div
                  className="mt-1 flex items-center rounded-md bg-white pl-3 outline outline-1 -outline-offset-1 outline-gray-300 focus-within:outline-2 focus-within:-outline-offset-2 focus-within:outline-indigo-600">
                <input
                    {...register('email')}
                    type="email"
                    autoComplete="email"
                    className="block min-w-0 grow bg-white py-1.5 pr-3 pl-1 text-base text-gray-900 placeholder:text-gray-400 focus:outline-none sm:text-sm/6"
                    placeholder="example@email.com"
                    disabled={isLoading}
                />
              </div>
              {errors.email && (
                <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="displayName" className="block text-sm font-medium text-gray-700">
                表示名
              </label>
              <div
                  className="mt-1 flex items-center rounded-md bg-white pl-3 outline outline-1 -outline-offset-1 outline-gray-300 focus-within:outline-2 focus-within:-outline-offset-2 focus-within:outline-indigo-600">
                <input
                    {...register('displayName')}
                    type="text"
                    autoComplete="name"
                    className="block min-w-0 grow bg-white py-1.5 pr-3 pl-1 text-base text-gray-900 placeholder:text-gray-400 focus:outline-none sm:text-sm/6"
                    placeholder="表示名を入力"
                    disabled={isLoading}
                />
              </div>
              {errors.displayName && (
                <p className="mt-1 text-sm text-red-600">{errors.displayName.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                パスワード
              </label>
              <div
                  className="mt-1 flex items-center rounded-md bg-white pl-3 outline outline-1 -outline-offset-1 outline-gray-300 focus-within:outline-2 focus-within:-outline-offset-2 focus-within:outline-indigo-600">
                <input
                    {...register('password')}
                    type="password"
                    autoComplete="new-password"
                    className="block min-w-0 grow bg-white py-1.5 pr-3 pl-1 text-base text-gray-900 placeholder:text-gray-400 focus:outline-none sm:text-sm/6"
                    placeholder="8文字以上、大文字・小文字・記号を含む"
                    disabled={isLoading}
                />
              </div>
              {errors.password && (
                <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700">
                パスワード確認
              </label>
              <div
                  className="mt-1 flex items-center rounded-md bg-white pl-3 outline outline-1 -outline-offset-1 outline-gray-300 focus-within:outline-2 focus-within:-outline-offset-2 focus-within:outline-indigo-600">
                <input
                    {...register('confirmPassword')}
                    type="password"
                    autoComplete="new-password"
                    className="block min-w-0 grow bg-white py-1.5 pr-3 pl-1 text-base text-gray-900 placeholder:text-gray-400 focus:outline-none sm:text-sm/6"
                    placeholder="パスワードを再入力"
                    disabled={isLoading}
                />
              </div>
              {errors.confirmPassword && (
                <p className="mt-1 text-sm text-red-600">{errors.confirmPassword.message}</p>
              )}
            </div>
          </div>

          {(error || errors.root) && (
            <div className="rounded-md bg-red-50 p-4">
              <div className="text-sm text-red-700">
                {error || errors.root?.message}
              </div>
            </div>
          )}

          <div>
            <button
              type="submit"
              disabled={isLoading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  アカウント作成中...
                </div>
              ) : (
                'アカウント作成'
              )}
            </button>
          </div>
        </form>
          </FadeIn>
        </Container>
      </Layout>
  );
};

export default SignupForm;