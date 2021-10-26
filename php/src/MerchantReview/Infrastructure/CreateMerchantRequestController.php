<?php

declare(strict_types=1);

namespace App\MerchantReview\Infrastructure;

use App\MerchantReview\Application\CreateMerchantReviewCommandHandler;
use Symfony\Component\HttpFoundation\Response;

final class CreateMerchantRequestController
{
    public function __construct(
        private CreateMerchantReviewCommandHandler $createMerchantReviewCommandHandler
    ) {}

    public function __invoke(): Response
    {
        // new CreateMerchantReviewCommandHandler(
        //     new CreateMerchantReviewCommand(

        //     ),
        // );

        return new Response();
    }
}
