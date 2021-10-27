<?php

declare(strict_types=1);

namespace App\MerchantReview\Application;

use App\MerchantReview\Domain\EventBus;
use App\MerchantReview\Domain\MerchantReview;
use App\MerchantReview\Domain\MerchantReviewCreated;
use App\MerchantReview\Domain\MerchantReviewRepository;
use App\MerchantReview\Domain\MerchantReviewUUIDGenerator;

final class CreateMerchantReviewCommandHandler
{
    public function __construct(
        private MerchantReviewUUIDGenerator $merchantReviewUUIDGenerator,
        private MerchantReviewRepository $merchantReviewRepository,
        private EventBus $eventBus,
    ) {}

    public function __invoke(
        CreateMerchantReviewCommand $command
    ): void {
        $review = MerchantReview::create(
            $this->merchantReviewUUIDGenerator->generate(),
            $command->getComment(),
            $command->getRating(),
        );

        $this->merchantReviewRepository->add($review);

        $this->eventBus->dispatch(

        );
    }
}
